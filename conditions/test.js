import { resolve_if } from './pkg/conditions.js'

var tests = {
  'true': {
    result: true,
    script: 'success',
    data: { 'success': true }
  },
  'false': {
    result: false,
    script: 'success',
    data: { 'success': false }
  },
  'test AND pass': {
    result: true,
    script: 'success && notFailed',
    data: { 'success': true, 'notFailed': true }
  },
  'test OR pass': {
    result: true,
    script: 'success || notFailed',
    data: { 'success': false, 'notFailed': true }
  },
  'test string pass': {
    result: true,
    script: 'name == "test"',
    data: { 'name': 'test' }
  },
  'test string fail': {
    result: false,
    script: 'name == "test"',
    data: { 'name': 'not test' }
  },
  'test string not match pass': {
    result: false,
    script: 'name != "test"',
    data: { 'name': 'test' }
  }
}

for (const testsKey in tests) {
  let test = tests[testsKey]

  let result = resolve_if(test.script, test.data)

  console.log('Test: ' + testsKey)
  console.log('Expected: ' + test.result)
  console.log('Actual: ' + result)

  if (test.result !== result) {
    throw 'Test failed'
  }
}