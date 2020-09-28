This folder contains the granted swagger schema, or schema properties that are not fit to be included in Terraform.

The granting is supposed to be done against each API version respectively, since API is expected to contain breaking changes between versions.

The folder structure is supposed to align with it in the Swagger repository, as the tool will lookup each grant list based on the relative path to the swagger spec.
