Transaction :
- Id : 1
- Amount : 10000
- Fee Fix : 500
- Fee Percetage : 5
- Net Amount ;
- Status : SETTLED
- CreatedAt : NOW


LOGIC CALLBACK : 9100

Wallet_transaction :
- Id : 1
- Type : Deposit
- WalletId : 2
- Source : Transaction
- SourceId : 1
- Amount : 9100
- Meta : {
  "TransactionId": 1
}

- Id : 2
- Type : Transfer
- WalletId : 3
- Source : Wallet
- SourceId : 1
- Amount : 9100
- Meta : {
  "WalletId": 1
}

- Id : 3
- Type : Deposit
- WalletId : 1
- Source : Transaction
- SourceId : 1
- Amount : 900
- Meta : {
  "WalletId": 1
}

Wallet :
- Id : 1
- UserId : 1
- Name : Cash
- Slug : cash
- Balance : 

- Id : 2
- UserId : 2
- Name : Pending
- Slug : pending
- Balance : 0 (Dikurangin)

- Id : 3
- UserId : 2
- Name : Cash
- Slug : cash
- Balance : 9100 (Ditambahin)
