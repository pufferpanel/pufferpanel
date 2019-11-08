package models

import "testing"

func TestUserView_Valid(t *testing.T) {
	type fields struct {
		Username string
		Email    string
		Password string
	}
	type args struct {
		allowEmpty bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "don't allow empty",
			fields: fields{
				Username: "",
				Email:    "",
				Password: "",
			},
			args:    args{allowEmpty: false},
			wantErr: true,
		},
		{
			name: "allow empty if requested",
			fields: fields{
				Username: "",
				Email:    "",
				Password: "",
			},
			args:    args{allowEmpty: true},
			wantErr: false,
		},
		{
			name: "test valid username",
			fields: fields{
				Username: "test1234",
				Email:    "",
				Password: "",
			},
			args:    args{allowEmpty: true},
			wantErr: false,
		},
		{
			name: "test invalid username",
			fields: fields{
				Username: "test & invalid",
				Email:    "",
				Password: "",
			},
			args:    args{allowEmpty: true},
			wantErr: true,
		},
		{
			name: "test invalid email",
			fields: fields{
				Username: "",
				Email:    "test",
				Password: "",
			},
			args:    args{allowEmpty: true},
			wantErr: true,
		},
		{
			name: "test valid email",
			fields: fields{
				Username: "",
				Email:    "test@example.com",
				Password: "",
			},
			args:    args{allowEmpty: true},
			wantErr: false,
		},
		{
			name: "test invalid email 2",
			fields: fields{
				Username: "",
				Email:    "test@com",
				Password: "",
			},
			args:    args{allowEmpty: true},
			wantErr: true,
		},
		{
			name: "test valid full object",
			fields: fields{
				Username: "validName",
				Email:    "valid@example.com",
				Password: "testing123!",
			},
			args:    args{allowEmpty: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &UserView{
				Username: tt.fields.Username,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			if err := model.Valid(tt.args.allowEmpty); (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
