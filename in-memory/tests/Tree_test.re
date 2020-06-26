open Jest;

module BasicTree = {
  type t = Tree.t(string);
  let rootID = "root-dir";

  let generator = FastCheck.array(FastCheck.string());

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

    FastCheck.expect(BasicTree.generator, checkAddNodes);
  })
});
