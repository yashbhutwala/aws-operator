package encryptionkey

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/aws-operator/service/controller/v23/controllercontext"
	"github.com/giantswarm/aws-operator/service/controller/v23/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.encrypter.EnsureCreatedEncryptionKey(ctx, customObject)
	if err != nil {
		return microerror.Mask(err)
	}

	encryptionKey, err := r.encrypter.EncryptionKey(ctx, customObject)
	if err != nil {
		return microerror.Mask(err)
	}

	cc, err := controllercontext.FromContext(ctx)
	if err != nil {
		return microerror.Mask(err)
	}
	cc.Status.Cluster.EncryptionKey = encryptionKey

	return nil
}
