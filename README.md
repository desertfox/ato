# **A**nsible Playbook **T**o **O**penAPIv3Schema
    Simple tool for creating OpenAPIv3Schema from Ansible playbook defaults/main.yml file. 
## Usage
* ```ato --crd=/path/to/crd.yaml```
## Notes
* Currently expects rigid file structure specific to my project layout. As such, it will look for the playbook 
by parsing out the names.singular value from the CustomResourceDefinition spec and using that to build the path.

* For example given the names.singular value of the CustomResourceDefinition is 'turtle', the bellow path
is where the application will look for yaml structure to map.  ```roles/turtle/defaults/main.yml```