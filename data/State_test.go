package data

import (
	"reflect"
	"sync"
	"testing"
)

func Test_Sha512(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"Test empty value returns expected.", "", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"Test value 'I like apples' returns expected. ", "I like apples", "921fde90136f6f6a0edf3cc4fced9533eff1133f945f9cbe2939d9c4d1135b65d5b2ce9bdf0d32934a0bd8e6302c9141a90a3888a1cc4887a553011008a5e994"},
		{"Test value with special characters returns expected. ", ":asdf#", "5ab74ab02ccf9d7617534419344415e93bc56622268ae781df98c344662bcc483e7818d31ed341964106942aa59a7fdd4aa267264229fb767682ec1ff4510729"},
		{"Test unicode character ă returns expected", "ă", "cd15b4ad6d85be0404acc0f921677ee6206893cee71642575660e8186a20d57b734c22e8aa57c331b7bf9e1fca5bb4f2b7f37d13f04ce9297717ee05fd82c6fc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hash(tt.value)
			if got != tt.want {
				t.Errorf("hash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test empty string returns expected.", args{value: ""}, ""},
		{"Test apple string returns expected.", args{value: "apple"}, "YXBwbGU="},
		{"Test hash string returns expected.",
			args{value: "cd15b4ad6d85be0404acc0f921677ee6206893cee71642575660e8186a20d57b734c22e8aa57c331b7bf9e1fca5bb4f2b7f37d13f04ce9297717ee05fd82c6fc"},
			"Y2QxNWI0YWQ2ZDg1YmUwNDA0YWNjMGY5MjE2NzdlZTYyMDY4OTNjZWU3MTY0MjU3NTY2MGU4MTg2YTIwZDU3YjczNGMyMmU4YWE1N2MzMzFiN2JmOWUxZmNhNWJiNGYyYjdmMzdkMTNmMDRjZTkyOTc3MTdlZTA1ZmQ4MmM2ZmM="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base64Encode(tt.args.value); got != tt.want {
				t.Errorf("base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		want *State
	}{
		{"test hashservice returns initial expectations", &State{
			hashCount:     0,
			hashCountLock: &sync.Mutex{},
			hashes:        &sync.Map{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashService_GetHashedPassword(t *testing.T) {
	type fields struct {
		hashCount     int64
		hashCountLock *sync.Mutex
		hashes        *sync.Map
	}
	type args struct {
		identifier int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue string
		wantFound bool
	}{
		{"Test no saved password returns not found.",
			fields{
				hashCount:     0,
				hashCountLock: &sync.Mutex{},
				hashes:        &sync.Map{},
			},
			args{
				identifier: 1,
			},
			"",
			false,
		},
		// TODO TW: Figure out how to preload the sync.Maps
		/*{"Test expected saved password returns found and expected value",
			fields{
				hashCount:     0,
				hashCountLock: &sync.Mutex{},
				hashes:        expectedSavedMap,
			},
			args{
				identifier: 1,
			},
			expectedSavedValue,
			true,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &State{
				hashCount:     tt.fields.hashCount,
				hashCountLock: tt.fields.hashCountLock,
				hashes:        tt.fields.hashes,
			}
			gotValue, gotFound := service.GetHashedPassword(tt.args.identifier)
			if gotFound != tt.wantFound {
				t.Errorf("GetHashedPassword() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
			if tt.wantFound && gotValue != tt.wantValue {
				t.Errorf("GetHashedPassword() gotValue = %v, want %v", gotValue, tt.wantValue)
			}

		})
	}
}

func TestState_SavePassword(t *testing.T) {
	type fields struct {
		hashCount     int64
		hashCountLock *sync.Mutex
		hashes        *sync.Map
	}
	type args struct {
		identifier int64
		password   string
	}
	type wantedArgs struct {
		identifier int64
		password   string
		found      bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantedArgs []wantedArgs
	}{
		{"Test Saving password uses the correct identifier for key and properly saves the password",
			fields{
				hashCount:     0,
				hashCountLock: &sync.Mutex{},
				hashes:        &sync.Map{},
			},
			args{
				identifier: 1,
				password:   "asdf",
			},
			[]wantedArgs{{
				identifier: 1,
				password:   base64Encode(hash("asdf")),
				found:      true,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &State{
				hashCount:     tt.fields.hashCount,
				hashCountLock: tt.fields.hashCountLock,
				hashes:        tt.fields.hashes,
			}
			state.SavePassword(tt.args.identifier, tt.args.password)
			for _, arg := range tt.wantedArgs {
				got, ok := tt.fields.hashes.Load(arg.identifier)
				if ok != arg.found {
					t.Errorf("SavePassword() = %v, want %v", ok, arg.found)
				}
				wanted := arg.password
				if got != wanted {
					t.Errorf("SavePassword() = %v, want %v", got, wanted)
				}
			}
		})
	}
}

func TestHashService_GetIdentifier(t *testing.T) {
	type fields struct {
		hashCount     int64
		hashCountLock *sync.Mutex
		hashes        *sync.Map
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Test no saved identifiers returns 1", fields{
			hashCount:     0,
			hashCountLock: &sync.Mutex{},
			hashes:        &sync.Map{},
		},
			1},
		{"Test one saved identifiers returns 2", fields{
			hashCount:     1,
			hashCountLock: &sync.Mutex{},
			hashes:        &sync.Map{},
		},
			2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &State{
				hashCount:     tt.fields.hashCount,
				hashCountLock: tt.fields.hashCountLock,
				hashes:        tt.fields.hashes,
			}
			if got := service.GetNextIdentifier(); got != tt.want {
				t.Errorf("GetNextIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
