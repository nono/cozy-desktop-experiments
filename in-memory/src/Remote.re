type dir = {
  id: Tree.id,
  parentID: Tree.id,
  name: string,
  rev: string,
};

type extra = {rev: string};

type tree = Tree.t(extra);

type changes = list(dir);

let rootID = "io.cozy.files.root-dir";

let emptyTree: tree = Tree.init(rootID);

let dirToNode = (dir: dir): Tree.node(extra) => {
  {
    id: dir.id,
    parentID: dir.parentID,
    name: dir.name,
    extra: {
      rev: dir.rev,
    },
  };
};
