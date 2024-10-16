package storage

import (
	"database/sql"
	"log"
)

// TransactionOperation
/** TransactionOperation godoc
Represents a function that has SQL instructions and will be executed inside a single transaction
*/
type TransactionOperation func() error

// UnitOfWork
/** UnitOfWork godoc
Represents a database transaction
- It is used to group database operations into a single transaction
- It is used to ensure that all operations are successful or none of them are */
type UnitOfWork struct {
	db          *sql.DB
	transaction *sql.Tx
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) Begin() error {
	tx, err := uow.db.Begin()
	if err != nil {
		return err
	}
	uow.transaction = tx
	log.Println("Transaction started...")
	return nil
}

func (uow *UnitOfWork) Commit() error {
	err := uow.transaction.Commit()
	if err != nil {
		return err
	}
	uow.transaction = nil
	log.Println("Transaction committed!")
	return nil
}

func (uow *UnitOfWork) Rollback() error {
	err := uow.transaction.Rollback()
	if err != nil {
		return err
	}
	uow.transaction = nil
	log.Println("Error in transaction. Rollback executed!")
	return nil
}

func (uow *UnitOfWork) ExecuteInTransaction(operation TransactionOperation) error {
	err := uow.Begin()
	if err != nil {
		return err
	}

	err = operation()
	if err != nil {
		_ = uow.Rollback()
		return err
	}

	err = uow.Commit()
	if err != nil {
		_ = uow.Rollback()
		return err
	}

	return nil
}
