// Generated by BUCKLESCRIPT, PLEASE EDIT WITH CARE
'use strict';

var Jest = require("@glennsl/bs-jest/src/jest.js");
var FastCheck = require("fast-check");

function expect(generator, predicate) {
  return Jest.Expect.toThrow(Jest.Expect.not_(Jest.Expect.expect(function (param) {
                      FastCheck.assert(FastCheck.property(generator, predicate));
                      
                    })));
}

exports.expect = expect;
/* Jest Not a pure module */