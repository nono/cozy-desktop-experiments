%raw
"const fc = require('fast-check');";

open Jest;

module BasicTree = {
  type t = Tree.t(string);
  let rootID = "root-dir";

  let fromArray = (ids: array(string)): t => {
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
    aux(Tree.init(rootID), List.fromArray(ids));
  };
};

describe("Tree", () => {
  test("addNode", () => {
    let checkAddNodes = ids => {
      let tree = BasicTree.fromArray(ids);
      let uniques = Set.String.fromArray(ids);
      Map.String.size(tree.nodes) == Set.String.size(uniques);
    };

    // XXX Don't let bucklescript remove the checkAddNodes function as an
    // optimization (as it don't see the call from the raw JS)
    Js.log(checkAddNodes);

    // TODO it would be nice to have some reason bindings for fast-check, but
    // calling it from raw JS is an acceptable work-around for the moment.
    Expect.(
      expect(() => {
        %raw
        "fc.assert(fc.property(fc.array(fc.string()), checkAddNodes))"
      })
      |> not
      |> toThrow
    );
  })
});
