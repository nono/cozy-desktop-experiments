open Remote

type t = {
  tree: tree,
  changes: changes,
}

let addDir = (cozy: t, dir: dir): t => {
  let tree = Remote.dirToNode(dir) |> Tree.addNode(cozy.tree)
  let changes = list{dir, ...cozy.changes}
  {tree: tree, changes: changes}
}
