package repository

import (
	"cryptocoin-server/model"
	"time"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

const server = "bolt://neo4j:password@localhost:7687"

// SaveTransaction saves transaction to the ledger.
func SaveTransaction(t *model.Transaction) error {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(server)
	if err != nil {
		return err
	}
	defer conn.Close()

	stmt, err := conn.PreparePipeline(
		"CREATE (n:Transaction {signature: {signature}, prevSignature: {prevSignature}, value: {value}, pubKey: {pubKey}, toAddress: {toAddress}, timestamp: {timestamp}})",
		"MATCH (c:Transaction),(p:Transaction) WHERE p.signature = {prevSignature} AND c.signature = {signature} CREATE (c)-[r:PREVIOUS]->(p)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecPipeline(
		map[string]interface{}{"signature": t.Signature, "prevSignature": t.PrevSignature, "value": t.Value, "pubKey": t.PubKey, "toAddress": t.ToAddress, "timestamp": t.Timestamp.Unix()},
		map[string]interface{}{"signature": t.Signature, "prevSignature": t.PrevSignature},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactions returns all latest transactions (limit 25).
func GetTransactions() ([]model.Transaction, error) {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := `
	MATCH
	  (n:Transaction)
	RETURN
	  n.signature, n.prevSignature, n.value, n.pubKey, n.toAddress, n.timestamp
	LIMIT {limit}`

	data, _, _, err := conn.QueryNeoAll(query, map[string]interface{}{"limit": 25})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	results := make([]model.Transaction, len(data))
	for i, row := range data {
		results[i] =
			model.Transaction{
				Signature:     row[0].(string),
				PrevSignature: row[1].(string),
				Value:         row[2].(int64),
				PubKey:        row[3].(string),
				ToAddress:     row[4].(string),
				Timestamp:     time.Unix(row[5].(int64), 0),
			}
	}

	return results, nil
}

// GetTransaction returns transaction by ID (Signature).
func GetTransaction(signature string, includeUsed bool) (*model.Transaction, error) {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := `
	MATCH
	  (n:Transaction)
	WHERE
	  n.signature = {signature}`

	if !includeUsed {
		query += `
		AND NOT ()-[:PREVIOUS]->(n) 
		`
	}

	query += `
	RETURN
	  n.signature, n.prevSignature, n.value, n.pubKey, n.toAddress, n.timestamp
	LIMIT {limit}`

	data, _, _, err := conn.QueryNeoAll(query, map[string]interface{}{"signature": signature, "limit": 1})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var transaction model.Transaction

	for _, row := range data {
		transaction =
			model.Transaction{
				Signature:     row[0].(string),
				PrevSignature: row[1].(string),
				Value:         row[2].(int64),
				PubKey:        row[3].(string),
				ToAddress:     row[4].(string),
				Timestamp:     time.Unix(row[5].(int64), 0),
			}
	}

	return &transaction, nil
}
