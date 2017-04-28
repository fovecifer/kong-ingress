package kong

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/client-go/pkg/api"
)

// HasKongFinalizer verify if the kong finalizer is set on the resource
func (d *Domain) HasKongFinalizer() bool {
	hasFinalizer := false
	for _, finalizer := range d.GetFinalizers() {
		if finalizer == Finalizer {
			hasFinalizer = true
			break
		}
	}
	return hasFinalizer
}

// IsMarkedForDeletion validates if the resource is set for deletion
func (d *Domain) IsMarkedForDeletion() bool {
	return d.DeletionTimestamp != nil || d.Status.DeletionTimestamp != nil
}

// IsPrimary validates if it's a primary domain
func (d *Domain) IsPrimary() bool {
	return len(d.Spec.Sub) == 0
}

// IsValidSharedDomain verifies if the shared domain it's a subdomain from the primary
func (d *Domain) IsValidSharedDomain() bool {
	return !d.IsPrimary() && d.IsValidDomain()
}

func (d *Domain) IsValidDomain() bool {
	if len(strings.Split(d.Spec.Sub, ".")) > 1 || len(strings.Split(d.Spec.PrimaryDomain, ".")) < 2 {
		return false
	}
	return true
}

func (d *Domain) GetDomain() string {
	if d.IsPrimary() {
		return d.GetPrimaryDomain()
	}
	return d.Spec.Sub + "." + d.Spec.PrimaryDomain
}

// GetDomainType returns the type of the resource: 'primary' or 'shared'
func (d *Domain) GetDomainType() string {
	if d.IsPrimary() {
		return "primary"
	}
	return "shared"
}

// GetPrimaryDomain returns the primary domain of the resource
func (d *Domain) GetPrimaryDomain() string {
	return d.Spec.PrimaryDomain
}

func (d *Domain) DeepCopy() (*Domain, error) {
	objCopy, err := api.Scheme.DeepCopy(d)
	if err != nil {
		return nil, err
	}
	copied, ok := objCopy.(*Domain)
	if !ok {
		return nil, fmt.Errorf("expected Domain, got %#v", objCopy)
	}
	return copied, nil
}

// IsUpdateExpired validates if the last update of the resource is expired
func (c *Domain) IsUpdateExpired(expireAfter time.Duration) bool {
	updatedAt := c.Status.LastUpdateTime.Add(expireAfter)
	if updatedAt.Before(time.Now().UTC()) {
		return true
	}
	return false
}

// // Required to satisfy Object interface
// func (d *Domainclaim) GetObjectKind() schema.ObjectKind {
// 	return &d.TypeMeta
// }

// // Required to satisfy ObjectMetaAccessor interface
// func (d *Domainclaim) GetObjectMeta() metav1.Object {
// 	return &d.ObjectMeta
// }

// // Required to satisfy Object interface
// func (dl *DomainclaimList) GetObjectKind() schema.ObjectKind {
// 	return &dl.TypeMeta
// }

// // Required to satisfy ListMetaAccessor interface
// func (dl *DomainclaimList) GetListMeta() metav1.List {
// 	return &dl.ListMeta
// }

// type DomainclaimListCopy DomainclaimList
// type DomainClaimCopy Domainclaim

// func (d *Domainclaim) UnmarshalJSON(data []byte) error {
// 	tmp := DomainClaimCopy{}
// 	err := json.Unmarshal(data, &tmp)
// 	if err != nil {
// 		return err
// 	}
// 	tmp2 := Domainclaim(tmp)
// 	*d = tmp2
// 	return nil
// }

// func (dl *DomainclaimList) UnmarshalJSON(data []byte) error {
// 	tmp := DomainclaimListCopy{}
// 	err := json.Unmarshal(data, &tmp)
// 	if err != nil {
// 		return err
// 	}
// 	tmp2 := DomainclaimList(tmp)
// 	*dl = tmp2
// 	return nil
// }
