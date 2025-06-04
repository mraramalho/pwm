package main

import (
	"fmt"
	"log"
)

type Record struct {
	User      string
	Password  string
	Encrypted bool
}

func NewRecord(user, password string) *Record {
	return &Record{user, password, false}
}

func (r *Record) encryptRecord() {
	if r.Encrypted {
		fmt.Println("Registro já está criptografado")
		return
	}

	u, err := encrypt(r.User)
	if err != nil {
		log.Println("Erro ao criptografar usuário", err)
		return
	}

	p, err := encrypt(r.Password)
	if err != nil {
		log.Println("Erro ao criptografar senha", err)
		return
	}

	r.User = u
	r.Password = p
	r.Encrypted = true
}

func (r *Record) decryptRecord() {
	if !r.Encrypted {
		fmt.Println("Registro não está criptografado")
		return
	}

	u, err := decrypt(r.User)
	if err != nil {
		log.Println("Erro ao descriptografar usuário", err)
		return
	}

	p, err := decrypt(r.Password)
	if err != nil {
		log.Println("Erro ao descriptografar senha", err)
		return
	}
	r.User = u
	r.Password = p
	r.Encrypted = false
}
