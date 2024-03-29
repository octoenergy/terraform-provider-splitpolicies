# Terraform Policy Splitter

This provider has one simple job - to take an array of documents and to combine them into documents no larger than a
given size. The driving force behind this is to take a list of AWS policies and combine them into large policy
documents that don't overflow then AWS-imposed limit of 6144 bytes. The chunk size can be configured though.

## Usage

```hcl
data "splitpolicies" "test" {
  policies = ["one", "two", "three"]
  maximum_chunk_size = 6
}

data "aws_iam_policy_document" "policy_docs" {
  for_each                = data.splitpolicies.test.chunks
  source_policy_documents = each.value
}

resource "aws_iam_policy" "policies" {
    for_each = data.aws_iam_policy_document.policy_docs
    policy   = each.value.json
}
```
