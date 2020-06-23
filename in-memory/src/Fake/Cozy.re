open Remote;

type t = {
  tree,
  changes,
};

let addDir = (cozy: t, dir: dir): t => {
  let tree = Remote.dirToNode(dir) |> Tree.addNode(cozy.tree);
  let changes = [dir, ...cozy.changes];
  {tree, changes};
};
