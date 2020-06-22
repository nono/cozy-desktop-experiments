type change = {
  id: string,
  name: string,
  parentID: string,
};

type changes =
  | List(change);
