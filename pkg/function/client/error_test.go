package client

import (
	"fmt"
	"reflect"
	"testing"
)

func TestErrKind_String(t *testing.T) {
	tests := []struct {
		name string
		ek   ErrKind
		want string
	}{
		{
			name: "retry",
			ek:   0,
			want: "Retryable",
		},
		{
			name: "non-retry",
			ek:   1,
			want: "NonRetryable",
		},
		{
			name: "unknown_1",
			ek:   2,
			want: "Unknown",
		},
		{
			name: "unknown_2",
			ek:   3,
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ek.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromError(t *testing.T) {
	retryable := Retryable
	nonRetryable := NonRetryable
	unknown := Unknown
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		args   args
		wantK  *ErrKind
		wantS  string
		wantOk bool
	}{
		{
			name: "good_retry",
			args: args{
				err: UDFError{
					ErrKind:    0,
					ErrMessage: "Retryable Error",
				},
			},
			wantK:  &retryable,
			wantS:  "Retryable Error",
			wantOk: true,
		},
		{
			name: "good_non_retry",
			args: args{
				err: UDFError{
					ErrKind:    1,
					ErrMessage: "NonRetryable Error",
				},
			},
			wantK:  &nonRetryable,
			wantS:  "NonRetryable Error",
			wantOk: true,
		},
		{
			name: "good_unknown_1",
			args: args{
				err: UDFError{
					ErrKind:    2,
					ErrMessage: "Unknown Error",
				},
			},
			wantK:  &unknown,
			wantS:  "Unknown Error",
			wantOk: true,
		},
		{
			name: "good_not_standard",
			args: args{
				err: fmt.Errorf("not a standard error"),
			},
			wantK:  &unknown,
			wantS:  "not a standard error",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotK, gotS, gotOk := FromError(tt.args.err)
			if !reflect.DeepEqual(gotK, tt.wantK) {
				t.Errorf("FromError() gotK = %v, want %v", gotK, tt.wantK)
			}
			if gotS != tt.wantS {
				t.Errorf("FromError() gotS = %v, want %v", gotS, tt.wantS)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FromError() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestUDFError_Error(t *testing.T) {
	type fields struct {
		ErrKind    ErrKind
		ErrMessage string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "retryable",
			fields: fields{
				ErrKind:    0,
				ErrMessage: "deadline exceeds",
			},
			want: "Retryable: deadline exceeds",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := UDFError{
				ErrKind:    tt.fields.ErrKind,
				ErrMessage: tt.fields.ErrMessage,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUDFError_ErrorKind(t *testing.T) {
	retryable := Retryable
	type fields struct {
		ErrKind    ErrKind
		ErrMessage string
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrKind
	}{
		{
			name: "retryable",
			fields: fields{
				ErrKind:    0,
				ErrMessage: "deadline exceeds",
			},
			want: &retryable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := UDFError{
				ErrKind:    tt.fields.ErrKind,
				ErrMessage: tt.fields.ErrMessage,
			}
			if got := e.ErrorKind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUDFError_ErrorMessage(t *testing.T) {
	type fields struct {
		ErrKind    ErrKind
		ErrMessage string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "retryable",
			fields: fields{
				ErrKind:    0,
				ErrMessage: "deadline exceeds",
			},
			want: "deadline exceeds",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := UDFError{
				ErrKind:    tt.fields.ErrKind,
				ErrMessage: tt.fields.ErrMessage,
			}
			if got := e.ErrorMessage(); got != tt.want {
				t.Errorf("ErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
