// Generated by BUCKLESCRIPT, PLEASE EDIT WITH CARE
'use strict';

var Jest = require("@glennsl/bs-jest/src/jest.js");
var Belt_List = require("bs-platform/lib/js/belt_List.js");
var Tree$InMemory = require("../src/Tree.bs.js");
var Belt_MapString = require("bs-platform/lib/js/belt_MapString.js");

const fc = require('fast-check');
;

Jest.describe("Tree", (function (param) {
        var buildTree = function (ids) {
          var _acc = Tree$InMemory.init("root-dir");
          var _ids = Belt_List.fromArray(ids);
          while(true) {
            var ids$1 = _ids;
            var acc = _acc;
            if (!ids$1) {
              return acc;
            }
            var id = ids$1.hd;
            _ids = ids$1.tl;
            _acc = Tree$InMemory.addNode(acc, {
                  id: /* ID */{
                    _0: id
                  },
                  parentID: /* ID */{
                    _0: acc.rootID
                  },
                  name: id,
                  extra: id
                });
            continue ;
          };
        };
        return Jest.test("addNode", (function (param) {
                      var checkAddNodes = function (ids) {
                        var tree = buildTree(ids);
                        return Belt_MapString.size(tree.nodes) === ids.length;
                      };
                      console.log(checkAddNodes);
                      return (fc.assert(fc.property(fc.array(fc.string()), checkAddNodes)));
                    }));
      }));

/*  Not a pure module */
