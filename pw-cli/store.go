package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

type Storage interface {
	GetRecord(domain string) (*Record, error)
	SaveRecord(domain string, record *Record) error
	UpdateRecord(domain string, record *Record) error
	DeleteRecord(domain string) error
	ListDomains() ([]string, error)
}

type Records struct {
	db *bbolt.DB
}

func New() (*Records, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appPath := filepath.Join(home, ".pw-cli")
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		err = os.Mkdir(appPath, 0700)
		if err != nil {
			return nil, err
		}
	}

	dbPath := filepath.Join(appPath, "my.db")
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("passwords"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Records{
		db: db,
	}, nil
}

func (r *Records) GetRecord(domain string) (*Record, error) {
	var record *Record

	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("passwords"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		v := bucket.Get([]byte(domain))
		if v == nil {
			return fmt.Errorf("%s not found", domain)
		}

		err := json.Unmarshal(v, &record)
		if err != nil {
			return err
		}
		record.decryptRecord()
		return nil
	})
	if err != nil {
		return record, err
	}
	return record, nil
}

func (r *Records) SaveRecord(domain string, record *Record) error {
	if saved_record, _ := r.GetRecord(domain); saved_record != nil {
		return fmt.Errorf("%s already exists", domain)
	}

	if !record.Encrypted {
		record.encryptRecord()
	}

	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("passwords"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		v, err := json.Marshal(record)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(domain), v)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Records) domainExists(domain string) bool {
	_, err := r.GetRecord(domain)
	return err == nil
}

func (r *Records) UpdateRecord(domain string, record *Record) error {
	if !r.domainExists(domain) {
		return fmt.Errorf("%s not found", domain)
	}
	if !record.Encrypted {
		record.encryptRecord()
	}
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("passwords"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		v, err := json.Marshal(record)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(domain), v)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Records) DeleteRecord(domain string) error {
	if !r.domainExists(domain) {
		return fmt.Errorf("%s not found", domain)
	}
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("passwords"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Delete([]byte(domain))
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Records) ListDomains() ([]string, error) {
	var domains []string
	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("passwords"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.ForEach(func(k, v []byte) error {
			domains = append(domains, string(k))
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return domains, nil
}
