package onepiecemessage_test

import (
	"github.com/straw-hat-team/onepiece/go/onepiece/eventsourcing/onepiecemessage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMessageType(t *testing.T) {
	type args struct {
		msgType string
	}
	tests := []struct {
		name    string
		args    args
		want    *onepiecemessage.MessageType
		wantErr error
	}{
		{
			name: "valid message type",
			args: args{msgType: "acmecorp.banking.bankaccount.v1.AccountOpened"},
			want: onepiecemessage.MessageType("acmecorp.banking.bankaccount.v1.AccountOpened").AsPtr(),
		},
		{
			name:    "missing tokens",
			args:    args{msgType: "bankaccount.v1.AccountOpened"},
			want:    nil,
			wantErr: onepiecemessage.ErrMessageTypeInvalid,
		},
		{
			name:    "invalid characters",
			args:    args{msgType: "bankaccount.v1.Account-Opened"},
			want:    nil,
			wantErr: onepiecemessage.ErrMessageTypeInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgType, err := onepiecemessage.NewMessageType(tt.args.msgType)
			require.Equal(t, tt.want, msgType)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
