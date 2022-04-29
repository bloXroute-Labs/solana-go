// Copyright 2022 bloXroute-Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"encoding/binary"
	"errors"
	"fmt"

	ag_solanago "github.com/bloXroute-Labs/solana-go"
	"github.com/bloXroute-Labs/solana-go/programs/serum"
	ag_format "github.com/bloXroute-Labs/solana-go/text/format"
	ag_binary "github.com/gagliardetto/binary"
	bin "github.com/gagliardetto/binary"
	ag_treeout "github.com/gagliardetto/treeout"
)

// CancelOrder (V2)
type CancelOrder struct {
	Side    *serum.Side
	OrderId *bin.Uint128

	// 0. `[writable]` market
	// 1. `[writable]` bids
	// 2. `[writable]` asks
	// 3. `[writable]` OpenOrders
	// 4. `[signer]` the OpenOrders owner
	// 5. `[writable]` event_q
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCancelOrderInstructionBuilder creates a new `CancelOrder` instruction builder.
func NewCancelOrderInstructionBuilder() *CancelOrder {
	nd := &CancelOrder{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 6),
	}
	return nd
}

func (c *CancelOrder) SetSide(s serum.Side) *CancelOrder {
	c.Side = &s
	return c
}

func (c *CancelOrder) SetOrderId(o bin.Uint128) *CancelOrder {
	c.OrderId = &o
	return c
}

func (c *CancelOrder) SetMarketAccount(a ag_solanago.PublicKey) *CancelOrder {
	c.AccountMetaSlice[0] = ag_solanago.Meta(a).WRITE()
	return c
}

func (c *CancelOrder) SetBidsAccount(a ag_solanago.PublicKey) *CancelOrder {
	c.AccountMetaSlice[1] = ag_solanago.Meta(a).WRITE()
	return c
}

func (c *CancelOrder) SetAsksAccount(a ag_solanago.PublicKey) *CancelOrder {
	c.AccountMetaSlice[2] = ag_solanago.Meta(a).WRITE()
	return c
}

func (c *CancelOrder) SetOpenOrdersAccount(a ag_solanago.PublicKey) *CancelOrder {
	c.AccountMetaSlice[3] = ag_solanago.Meta(a).WRITE()
	return c
}

func (c *CancelOrder) SetOwnerAccount(a ag_solanago.PublicKey) *CancelOrder {
	c.AccountMetaSlice[4] = ag_solanago.Meta(a).WRITE().SIGNER()
	return c
}

func (c *CancelOrder) SetEventQueueAccount(a ag_solanago.PublicKey) *CancelOrder {
	c.AccountMetaSlice[5] = ag_solanago.Meta(a).WRITE()
	return c
}

func (c *CancelOrder) MarketAccount() *ag_solanago.AccountMeta {
	return c.AccountMetaSlice[0]
}

func (c *CancelOrder) BidsAccount() *ag_solanago.AccountMeta {
	return c.AccountMetaSlice[1]
}

func (c *CancelOrder) AsksAccount() *ag_solanago.AccountMeta {
	return c.AccountMetaSlice[2]
}

func (c *CancelOrder) OpenOrdersAccount() *ag_solanago.AccountMeta {
	return c.AccountMetaSlice[3]
}

func (c *CancelOrder) OwnerAccount() *ag_solanago.AccountMeta {
	return c.AccountMetaSlice[4]
}

func (c *CancelOrder) EventQueueAccount() *ag_solanago.AccountMeta {
	return c.AccountMetaSlice[5]
}

func (c CancelOrder) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   c,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_CancelOrder, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (c CancelOrder) ValidateAndBuild() (*Instruction, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c.Build(), nil
}

func (c *CancelOrder) Validate() error {
	if c.Side == nil {
		return errors.New("CancelOrder.Side parameter is not set")
	}
	if c.OrderId == nil {
		return errors.New("CancelOrder.OrderId parameter is not set")
	}

	for accIndex, acc := range c.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("CancelOrder.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}

func (c *CancelOrder) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).ParentFunc(func(programBranch ag_treeout.Branches) {
		programBranch.Child(ag_format.Instruction("CancelOrder")).
			ParentFunc(func(instructionBranch ag_treeout.Branches) {
				// Parameters of the instruction:
				instructionBranch.Child("Params").ParentFunc(func(parms ag_treeout.Branches) {
					parms.Child(ag_format.Param("    Side", *c.Side))
					parms.Child(ag_format.Param(" OrderId", *c.OrderId))
				})
				// Accounts of the instruction:
				instructionBranch.Child("Accounts").ParentFunc(func(accs ag_treeout.Branches) {
					accs.Child(ag_format.Meta("     Market", c.AccountMetaSlice[0]))
					accs.Child(ag_format.Meta("       Bids", c.AccountMetaSlice[1]))
					accs.Child(ag_format.Meta("       Asks", c.AccountMetaSlice[2]))
					accs.Child(ag_format.Meta(" OpenOrders", c.AccountMetaSlice[3]))
					accs.Child(ag_format.Meta("      Owner", c.AccountMetaSlice[4]))
					accs.Child(ag_format.Meta("  EvenQueue", c.AccountMetaSlice[5]))
				})
			})
	})
}

func (c CancelOrder) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	err := encoder.Encode(*c.Side)
	if err != nil {
		return err
	}
	return encoder.Encode(*c.OrderId)
}

func (c *CancelOrder) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	err := decoder.Decode(&c.Side)
	if err != nil {
		return err
	}
	return decoder.Decode(&c.OrderId)
}

// NewCancelOrderInstruction declares a new CancelOrder instruction with the provided parameters and accounts.
func NewCancelOrderInstruction(
	// Parameters:
	side serum.Side,
	orderId bin.Uint128,
	// Accounts:
	marketAccount ag_solanago.PublicKey,
	bidsAccount ag_solanago.PublicKey,
	asksAccount ag_solanago.PublicKey,
	openOrdersAccount ag_solanago.PublicKey,
	ownerAccount ag_solanago.PublicKey,
	eventQeueueAccount ag_solanago.PublicKey) *CancelOrder {
	return NewCancelOrderInstructionBuilder().
		SetSide(side).
		SetOrderId(orderId).
		SetMarketAccount(marketAccount).
		SetBidsAccount(marketAccount).
		SetAsksAccount(marketAccount).
		SetOpenOrdersAccount(openOrdersAccount).
		SetOwnerAccount(ownerAccount).
		SetEventQueueAccount(ownerAccount)
}
