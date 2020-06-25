%raw
"const fc = require('fast-check');";

open Jest;

describe("Tree", () => {
  let buildTree = (ids: array(string)): Tree.t(string) => {
    let rec aux = (acc, ids) => {
      switch (ids) {
      | [] => acc
      | [id, ...rest] =>
        aux(
          Tree.addNode(
            acc,
            {
              id: Tree.ID(id),
              parentID: Tree.ID(acc.rootID),
              name: id,
              extra: id,
            },
          ),
          rest,
        )
      };
    };
    aux(Tree.init("root-dir"), List.fromArray(ids));
  };

  test("addNode", () => {
    let checkAddNodes = ids => {
      let tree = buildTree(ids);
      Map.String.size(tree.nodes) == Array.size(ids);
    };

    // XXX Don't let bucklescript remove the checkAddNodes function as an
    // optimization (as it don't see the call from the raw JS)
    Js.log(checkAddNodes);

    // TODO it would be nice to have some reason bindings for fast-check, but
    // calling it from raw JS is an acceptable work-around for the moment
    %raw
    "fc.assert(fc.property(fc.array(fc.string()), checkAddNodes))";
  });
});
