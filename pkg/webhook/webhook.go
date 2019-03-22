package webhook

import (
	kuberule "github.com/chickenzord/kube-rule/pkg/apis/kuberule/v1alpha1"
	"github.com/chickenzord/kube-rule/pkg/config"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

func createMutatePodsWebhook(mgr manager.Manager) (*admission.Webhook, error) {
	return builder.NewWebhookBuilder().
		Name("mutatepods.kuberule.chickenzord.com").
		Mutating().
		Operations(
			admissionregistrationv1beta1.Create,
			admissionregistrationv1beta1.Update,
		).
		ForType(&corev1.Pod{}).
		Handlers(&podMutationHandler{
			client:  mgr.GetClient(),
			decoder: mgr.GetAdmissionDecoder(),
		}).
		FailurePolicy(admissionregistrationv1beta1.Ignore).
		WithManager(mgr).
		Build()
}

func createValidatePodRulesWebhook(mgr manager.Manager) (*admission.Webhook, error) {
	return builder.NewWebhookBuilder().
		Name("validatepodrules.kuberule.chickenzord.com").
		Validating().
		Operations(
			admissionregistrationv1beta1.Create,
			admissionregistrationv1beta1.Update,
		).
		ForType(&kuberule.PodRule{}).
		Handlers(&podRuleValidationHandler{
			client:  mgr.GetClient(),
			decoder: mgr.GetAdmissionDecoder(),
		}).
		FailurePolicy(admissionregistrationv1beta1.Ignore).
		WithManager(mgr).
		Build()
}

func createMutatePodRulesWebhook(mgr manager.Manager) (*admission.Webhook, error) {
	return builder.NewWebhookBuilder().
		Name("mutatepodrules.kuberule.chickenzord.com").
		Mutating().
		Operations(
			admissionregistrationv1beta1.Create,
			admissionregistrationv1beta1.Update,
		).
		ForType(&kuberule.PodRule{}).
		Handlers(&podRuleMutationHandler{
			client:  mgr.GetClient(),
			decoder: mgr.GetAdmissionDecoder(),
		}).
		FailurePolicy(admissionregistrationv1beta1.Ignore).
		WithManager(mgr).
		Build()
}

func createServer(mgr manager.Manager) (*webhook.Server, error) {
	return webhook.NewServer(config.AppName, mgr, webhook.ServerOptions{
		CertDir: config.CertDir,
		BootstrapOptions: &webhook.BootstrapOptions{
			MutatingWebhookConfigName:   config.AppName,
			ValidatingWebhookConfigName: config.AppName,

			Secret: &types.NamespacedName{
				Namespace: config.Namespace,
				Name:      config.SecretName,
			},

			Service: &webhook.Service{
				Namespace: config.Namespace,
				Name:      config.ServiceName,
				// Selectors should select the pods that runs this webhook server.
				Selectors: config.ServiceSelector,
			},
		},
	})
}

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs = []func(manager.Manager) error{
	func(mgr manager.Manager) error {
		mutatePodsWebhook, err := createMutatePodsWebhook(mgr)
		if err != nil {
			return err
		}

		validatePodRulesWebhook, err := createValidatePodRulesWebhook(mgr)
		if err != nil {
			return err
		}

		mutatePodRulesWebhook, err := createMutatePodRulesWebhook(mgr)
		if err != nil {
			return err
		}

		server, err := createServer(mgr)
		if err != nil {
			return err
		}

		return server.Register(
			mutatePodsWebhook,
			validatePodRulesWebhook,
			mutatePodRulesWebhook,
		)
	},
}

// AddToManager adds all Controllers to the Manager
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations;validatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
func AddToManager(m manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
