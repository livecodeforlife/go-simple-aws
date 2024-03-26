package infra

import (
	"reflect"
	"testing"
)

type VpcCreatorImpl struct {
	handle VPCHandle
	calls  uint
}

func (v *VpcCreatorImpl) create(vpc VPCDefinition) (*VPCHandle, error) {
	v.calls++
	return &v.handle, nil
}

func Test_createOrUpdateVpc(t *testing.T) {
	type args struct {
		def VPCDefinition
		vpc VPCCreator
	}

	vpcCreator := VpcCreatorImpl{
		handle: VPCHandle{},
		calls:  0,
	}

	tests := []struct {
		name    string
		args    args
		want    *VPCHandle
		wantErr bool
	}{
		{
			name: "Given a definition and a vpc creator, the function should call the create method",
			args: args{
				def: VPCDefinition{},
				vpc: &vpcCreator,
			},
			wantErr: false,
			want:    &vpcCreator.handle,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createOrUpdateVpc(tt.args.def, tt.args.vpc)
			if (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateVpc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createOrUpdateVpc() = %v, want %v", got, tt.want)
			}
		})
	}
}
