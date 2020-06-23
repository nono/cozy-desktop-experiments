// Generated by BUCKLESCRIPT, PLEASE EDIT WITH CARE
'use strict';

var Belt_MapString = require("bs-platform/lib/js/belt_MapString.js");

function init(rootID) {
  return {
          rootID: rootID,
          nodes: undefined
        };
}

function addNode(tree, node) {
  var id = node.id;
  var nodes = Belt_MapString.set(tree.nodes, id._0, node);
  return {
          rootID: tree.rootID,
          nodes: nodes
        };
}

exports.init = init;
exports.addNode = addNode;
/* No side effect */