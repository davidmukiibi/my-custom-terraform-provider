# Custom Terraform Provider making API calls to a custom API endpoint

## To make use of the newly built provider, follow the steps below

1. While in the root folder, run `make install`. this will:
   1. Build the Go binary.
   2. Create the folders for the new provider in the same location terraform stores its plugins.
   3. and then move the built binary to that folder.
2. run `export URL=$(minikube service --url stage-backend-service -n stage)` to set the environment variable for the API url and port number.
3. cd into the `examples` folder where there's a `main.tf` file.
   1. run `terraform init` to initialize the newly built terraform provider.
   2. run `terraform apply --auto-approve` to create resources declared in the `main.tf` file.
4. And viola! You successfully deployed a custom Terraform Provider calling a custom API.
5. A more detailed article on how the moving pieces were put together will come soon on medium [here](https://medium.com/@david.mukiibiq)

This build was inspired from the Harshicorp's `Terraform` who have a good example (hashicups) and documentation [here](https://www.hashicorp.com/blog/writing-custom-terraform-providers)
