type id =
  | ID(string);

type node('a) = {
  id,
  parentID: id,
  name: string,
  extra: 'a,
};

type t('a) = {
  rootID: string,
  nodes: Map.String.t(node('a)),
};

let init = (rootID: string): t('a) => {
  {rootID, nodes: Map.String.empty};
};

// TODO check that we don't already have a node with the same id
// TODO check that we don't already have a node with the same parentID+name
let addNode = (tree: t('a), node: node('a)): t('a) => {
  let ID(id) = node.id;
  let nodes = Map.String.set(tree.nodes, id, node);
  {rootID: tree.rootID, nodes};
};
