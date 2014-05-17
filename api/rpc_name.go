package api

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gitchain/gitchain/env"
	"github.com/gitchain/gitchain/server"
	"github.com/gitchain/gitchain/transaction"
)

type NameService struct{}

type NameReservationArgs struct {
	Alias string
	Name  string
}

type NameReservationReply struct {
	Id     string
	Random string
}

func (srv *NameService) NameReservation(r *http.Request, args *NameReservationArgs, reply *NameReservationReply) error {
	key, err := env.DB.GetKey(args.Alias)
	if err != nil {
		return err
	}
	if key == nil {
		return errors.New("can't find the key")
	}
	tx, random := transaction.NewNameReservation(args.Name, &key.PublicKey)
	reply.Id = hex.EncodeToString(tx.Hash())
	reply.Random = hex.EncodeToString(random)
	server.BroadcastTransaction(tx)
	return nil
}

type NameAllocationArgs struct {
	Alias  string
	Name   string
	Random string
}

type NameAllocationReply struct {
	Id string
}

func (srv *NameService) NameAllocation(r *http.Request, args *NameAllocationArgs, reply *NameAllocationReply) error {
	key, err := env.DB.GetKey(args.Alias)
	if err != nil {
		return err
	}
	if key == nil {
		return errors.New("can't find the key")
	}
	random, err := hex.DecodeString(args.Random)
	if err != nil {
		return err
	}
	tx, err := transaction.NewNameAllocation(args.Name, random, key)
	if err != nil {
		return err
	}

	reply.Id = hex.EncodeToString(tx.Hash())
	server.BroadcastTransaction(tx)
	return nil
}
