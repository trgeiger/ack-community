// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package cache_subnet_group

import (
	"context"

	ackv1alpha1 "github.com/aws/aws-controllers-k8s/apis/core/v1alpha1"
	ackcompare "github.com/aws/aws-controllers-k8s/pkg/compare"
	ackerr "github.com/aws/aws-controllers-k8s/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/elasticache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws/aws-controllers-k8s/services/elasticache/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = &aws.JSONValue{}
	_ = &svcsdk.ElastiCache{}
	_ = &svcapitypes.CacheSubnetGroup{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	input, err := rm.newListRequestPayload(r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.DescribeCacheSubnetGroupsWithContext(ctx, input)
	if respErr != nil {
		if awsErr, ok := ackerr.AWSError(respErr); ok && awsErr.Code() == "CacheSubnetGroupNotFoundFault" {
			return nil, ackerr.NotFound
		}
		return nil, respErr
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if len(resp.CacheSubnetGroups) == 0 {
		return nil, ackerr.NotFound
	}
	found := false
	for _, elem := range resp.CacheSubnetGroups {
		if elem.ARN != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.ARN)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.CacheSubnetGroupDescription != nil {
			ko.Spec.CacheSubnetGroupDescription = elem.CacheSubnetGroupDescription
		}
		if elem.CacheSubnetGroupName != nil {
			ko.Spec.CacheSubnetGroupName = elem.CacheSubnetGroupName
		}
		if elem.Subnets != nil {
			f3 := []*svcapitypes.Subnet{}
			for _, f3iter := range elem.Subnets {
				f3elem := &svcapitypes.Subnet{}
				if f3iter.SubnetAvailabilityZone != nil {
					f3elemf0 := &svcapitypes.AvailabilityZone{}
					if f3iter.SubnetAvailabilityZone.Name != nil {
						f3elemf0.Name = f3iter.SubnetAvailabilityZone.Name
					}
					f3elem.SubnetAvailabilityZone = f3elemf0
				}
				if f3iter.SubnetIdentifier != nil {
					f3elem.SubnetIdentifier = f3iter.SubnetIdentifier
				}
				f3 = append(f3, f3elem)
			}
			ko.Status.Subnets = f3
		}
		if elem.VpcId != nil {
			ko.Status.VPCID = elem.VpcId
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}

	return &resource{ko}, nil
}

// newListRequestPayload returns SDK-specific struct for the HTTP request
// payload of the List API call for the resource
func (rm *resourceManager) newListRequestPayload(
	r *resource,
) (*svcsdk.DescribeCacheSubnetGroupsInput, error) {
	res := &svcsdk.DescribeCacheSubnetGroupsInput{}

	if r.ko.Spec.CacheSubnetGroupName != nil {
		res.SetCacheSubnetGroupName(*r.ko.Spec.CacheSubnetGroupName)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a new resource with any fields in the Status field filled in
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	input, err := rm.newCreateRequestPayload(r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.CreateCacheSubnetGroupWithContext(ctx, input)
	if respErr != nil {
		return nil, respErr
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.CacheSubnetGroup.Subnets != nil {
		f3 := []*svcapitypes.Subnet{}
		for _, f3iter := range resp.CacheSubnetGroup.Subnets {
			f3elem := &svcapitypes.Subnet{}
			if f3iter.SubnetAvailabilityZone != nil {
				f3elemf0 := &svcapitypes.AvailabilityZone{}
				if f3iter.SubnetAvailabilityZone.Name != nil {
					f3elemf0.Name = f3iter.SubnetAvailabilityZone.Name
				}
				f3elem.SubnetAvailabilityZone = f3elemf0
			}
			if f3iter.SubnetIdentifier != nil {
				f3elem.SubnetIdentifier = f3iter.SubnetIdentifier
			}
			f3 = append(f3, f3elem)
		}
		ko.Status.Subnets = f3
	}
	if resp.CacheSubnetGroup.VpcId != nil {
		ko.Status.VPCID = resp.CacheSubnetGroup.VpcId
	}

	ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{OwnerAccountID: &rm.awsAccountID}
	ko.Status.Conditions = []*ackv1alpha1.Condition{}
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	r *resource,
) (*svcsdk.CreateCacheSubnetGroupInput, error) {
	res := &svcsdk.CreateCacheSubnetGroupInput{}

	if r.ko.Spec.CacheSubnetGroupDescription != nil {
		res.SetCacheSubnetGroupDescription(*r.ko.Spec.CacheSubnetGroupDescription)
	}
	if r.ko.Spec.CacheSubnetGroupName != nil {
		res.SetCacheSubnetGroupName(*r.ko.Spec.CacheSubnetGroupName)
	}
	if r.ko.Spec.SubnetIDs != nil {
		f2 := []*string{}
		for _, f2iter := range r.ko.Spec.SubnetIDs {
			var f2elem string
			f2elem = *f2iter
			f2 = append(f2, &f2elem)
		}
		res.SetSubnetIds(f2)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	r *resource,
	diffReporter *ackcompare.Reporter,
) (*resource, error) {
	input, err := rm.newUpdateRequestPayload(r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.ModifyCacheSubnetGroupWithContext(ctx, input)
	if respErr != nil {
		return nil, respErr
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.CacheSubnetGroup.Subnets != nil {
		f3 := []*svcapitypes.Subnet{}
		for _, f3iter := range resp.CacheSubnetGroup.Subnets {
			f3elem := &svcapitypes.Subnet{}
			if f3iter.SubnetAvailabilityZone != nil {
				f3elemf0 := &svcapitypes.AvailabilityZone{}
				if f3iter.SubnetAvailabilityZone.Name != nil {
					f3elemf0.Name = f3iter.SubnetAvailabilityZone.Name
				}
				f3elem.SubnetAvailabilityZone = f3elemf0
			}
			if f3iter.SubnetIdentifier != nil {
				f3elem.SubnetIdentifier = f3iter.SubnetIdentifier
			}
			f3 = append(f3, f3elem)
		}
		ko.Status.Subnets = f3
	}
	if resp.CacheSubnetGroup.VpcId != nil {
		ko.Status.VPCID = resp.CacheSubnetGroup.VpcId
	}

	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	r *resource,
) (*svcsdk.ModifyCacheSubnetGroupInput, error) {
	res := &svcsdk.ModifyCacheSubnetGroupInput{}

	if r.ko.Spec.CacheSubnetGroupDescription != nil {
		res.SetCacheSubnetGroupDescription(*r.ko.Spec.CacheSubnetGroupDescription)
	}
	if r.ko.Spec.CacheSubnetGroupName != nil {
		res.SetCacheSubnetGroupName(*r.ko.Spec.CacheSubnetGroupName)
	}
	if r.ko.Spec.SubnetIDs != nil {
		f2 := []*string{}
		for _, f2iter := range r.ko.Spec.SubnetIDs {
			var f2elem string
			f2elem = *f2iter
			f2 = append(f2, &f2elem)
		}
		res.SetSubnetIds(f2)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) error {
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return err
	}
	_, respErr := rm.sdkapi.DeleteCacheSubnetGroupWithContext(ctx, input)
	return respErr
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteCacheSubnetGroupInput, error) {
	res := &svcsdk.DeleteCacheSubnetGroupInput{}

	if r.ko.Spec.CacheSubnetGroupName != nil {
		res.SetCacheSubnetGroupName(*r.ko.Spec.CacheSubnetGroupName)
	}

	return res, nil
}