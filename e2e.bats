#!/usr/bin/env bats

@test "reject nodeport service when nodeport is disabled" {
  run kwctl run annotated-policy.wasm -r test_data/service-nodeport.json --settings-json '{"disable_nodeport": true}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
  [ $(expr "$output" : ".*Service 'test-service' of type NodePort is not allowed as per policy configuration.*") -ne 0 ]
}

@test "accept nodeport service when nodeport is allowed" {
  run kwctl run annotated-policy.wasm -r test_data/service-nodeport.json --settings-json '{"disable_nodeport": false}'
  
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "accept clusterip service when nodeport is disabled" {
  run kwctl run annotated-policy.wasm -r test_data/service-clusterip.json --settings-json '{"disable_nodeport": true}'
  
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "accept service with default settings" {
  run kwctl run annotated-policy.wasm -r test_data/service-nodeport.json
  
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}