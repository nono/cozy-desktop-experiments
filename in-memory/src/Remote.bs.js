// Generated by BUCKLESCRIPT, PLEASE EDIT WITH CARE
'use strict';

var Tree$InMemory = require("./Tree.bs.js");

var rootID = "io.cozy.files.root-dir";

var emptyTree = Tree$InMemory.init(rootID);

function dirToNode(dir) {
  return {
          id: dir.id,
          parentID: dir.parentID,
          name: dir.name,
          extra: {
            rev: dir.rev
          }
        };
}

exports.rootID = rootID;
exports.emptyTree = emptyTree;
exports.dirToNode = dirToNode;
/* emptyTree Not a pure module */