open Model;

type handler = model => (model, effects);

let incrementTicked: handler =
  current => {
    ({...current, ticked: current.ticked + 1}, []);
  };

let checkFetchChanges: handler =
  current => {
    switch (current.changesState) {
    | ChangesFeedNeverFetched => (
        {...current, changesState: ChangesFeedCurrentlyFetching},
        [FetchChangesFeed],
      )
    | _ => (current, [])
    };
  };

let combineHandlers = (handlers: list(handler)): handler => {
  let rec aux = (handlers: list(handler), current: model, acc: effects) => {
    switch (handlers) {
    | [] => (current, acc)
    | [handler, ...rest] =>
      let (next, effects) = handler(current);
      aux(rest, next, List.concat(acc, effects));
    };
  };
  model => aux(handlers, model, []);
};

let handleTick = combineHandlers([incrementTicked, checkFetchChanges]);

let update = (current: model, event: event): (model, effects) => {
  switch (event) {
  | Tick => handleTick(current)
  | _ => (current, [])
  };
};
