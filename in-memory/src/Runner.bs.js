// Generated by BUCKLESCRIPT, PLEASE EDIT WITH CARE
'use strict';

var Sync$InMemory = require("./Sync.bs.js");
var Model$InMemory = require("./Model.bs.js");

function apply(effects) {
  console.log([
        "apply",
        effects
      ]);
  
}

function run(config) {
  var initial = Model$InMemory.init(config);
  var model = {
    contents: initial
  };
  return setInterval((function (param) {
                var $$event = /* Tick */0;
                var match = Sync$InMemory.update(model.contents, $$event);
                model.contents = match[0];
                return apply(match[1]);
              }), 1);
}

exports.apply = apply;
exports.run = run;
/* Sync-InMemory Not a pure module */
