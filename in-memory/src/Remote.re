type id = string;

type node = {
  id,
  name: string,
  parentID: string,
};

type tree = Tree.t(node);

type change = node;

type changes =
  | List(change);

let emptyTree: tree = {
  rootID: "io.cozy.files.root-dir",
  nodes: Map.String.empty,
};
