package schemas

import (
	"github.com/ServiceComb/go-chassis/examples/schemas/employ"

	"golang.org/x/net/context"
)

//EmployServer is a struct
type EmployServer struct{}

//AddEmploy 这里实现服务端接口中的方法。
func (s *EmployServer) AddEmploy(ctx context.Context, in *employ.EmployRequest) (*employ.EmployResponse, error) {
	in.EmployList = append(in.EmployList, in.Employ)
	return &employ.EmployResponse{Employ: nil, EmployList: in.EmployList}, nil
}

//EditEmploy is a method used to edit employ
func (s *EmployServer) EditEmploy(ctx context.Context, in *employ.EmployRequest) (*employ.EmployResponse, error) {
	for index, value := range in.EmployList {
		if in.Name == value.Name {
			in.EmployList[index] = in.EmployList[len(in.EmployList)-1]
			tempList := in.EmployList[:len(in.EmployList)-1]
			tempList = append(tempList, in.Employ)
			return &employ.EmployResponse{Employ: nil, EmployList: tempList}, nil
		}
	}
	return &employ.EmployResponse{Employ: nil, EmployList: in.EmployList}, nil
}

//GetEmploys is a method used to get employs
func (s *EmployServer) GetEmploys(ctx context.Context, in *employ.EmployRequest) (*employ.EmployResponse, error) {
	for index, value := range in.EmployList {
		if in.Name == value.Name {
			temp := in.EmployList[index]
			return &employ.EmployResponse{Employ: temp, EmployList: in.EmployList}, nil
		}
	}
	return &employ.EmployResponse{Employ: nil, EmployList: in.EmployList}, nil
}

//DeleteEmploys is a method used to delete employ
func (s *EmployServer) DeleteEmploys(ctx context.Context, in *employ.EmployRequest) (*employ.EmployResponse, error) {
	for index, value := range in.EmployList {
		if in.Name == value.Name {
			in.EmployList[index] = in.EmployList[len(in.EmployList)-1]
			tempList := in.EmployList[:len(in.EmployList)-1]
			return &employ.EmployResponse{Employ: nil, EmployList: tempList}, nil
		}
	}
	return &employ.EmployResponse{Employ: nil, EmployList: in.EmployList}, nil
}
