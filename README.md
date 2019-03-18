# deploy wizard

## Stories

### [Kubernetes Deployment Wizard](https://rally1.rallydev.com/#/281494669804d/detail/userstory/294802121048?fdp=true)

Build a wizard that allows a team to create a fully fledged deployment using Kubernetes and Kustomize but with the need of manually writing YAML files.

The wizard will drive the application team to specify the parameters that are important to them and guide them through deployments in Kubernetes.

The wizard reads the YAML template and display a UI where the variables are made accessible and ready to be filled in. When the user has completed all the actions the Wizard will create the Kustomize set of YAML in the Git repo, folder and target branch indicated.

Template Definition:
https://fusion.mastercard.int/confluence/display/DNA/YAML+Templates+for+GitOps+Deployments

Wireframes:
https://fusion.mastercard.int/confluence/display/DNA/Kubernetes+Application+Deployment+Wizard

#### Definition of Done

As a user I will be able to access a Wizard UI that allows me to create a deployment template for my application. Once I finish the process and I hit "Create Deployment" I will see a set of Kustomize deployment files created in my git repo. Using either the Kustomize command line or ArgoCD I should be able to create my deployment in my namespace.

### [Kubernetes Deployment Wizard Wireframes](https://rally1.rallydev.com/#/281494669804d/detail/userstory/294814394640?fdp=true)

Build a set of wireframes that illustrate how the Wizard should work and give a base for the UI development.

Wireframes:
https://fusion.mastercard.int/confluence/display/DNA/Kubernetes+Application+Deployment+Wizard

#### Definition of Done

Have a full set of wireframes illustrating the visual representation as well as the user interaction.
