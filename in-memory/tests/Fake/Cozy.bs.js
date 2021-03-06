// Generated by BUCKLESCRIPT, PLEASE EDIT WITH CARE
'use strict';

var Tree$InMemory = require("../../src/Tree.bs.js");
var Remote$InMemory = require("../../src/Remote.bs.js");

function addDir(cozy, dir) {
  var tree = Tree$InMemory.addNode(cozy.tree, Remote$InMemory.dirToNode(dir));
  var changes_1 = cozy.changes;
  var changes = {
    hd: dir,
    tl: changes_1
  };
  return {
          tree: tree,
          changes: changes
        };
}

exports.addDir = addDir;
/* Remote-InMemory Not a pure module */
