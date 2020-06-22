open Model;

type handler = t => (t, effects);

let incrementTicked: handler =
  current => {
    ({...current, ticked: current.ticked + 1}, []);
  };

// TODO: fetches again the changes feed after some time
let checkFetchChanges: handler =
  model => {
    switch (model.states.changes) {
    | ChangesFeedNeverFetched => (
        {
          ...model,
          states: {
            changes: ChangesFeedCurrentlyFetching,
          },
        },
        [FetchChangesFeed],
      )
    | _ => (model, [])
    };
  };

let combineHandlers = (handlers: list(handler)): handler => {
  let rec aux = (handlers: list(handler), model: t, acc: effects) => {
    switch (handlers) {
    | [] => (model, acc)
    | [handler, ...rest] =>
      let (next, effects) = handler(model);
      aux(rest, next, List.concat(acc, effects));
    };
  };
  model => aux(handlers, model, []);
};

let handleTick = combineHandlers([incrementTicked, checkFetchChanges]);

let update = (model: t, event: event): (t, effects) => {
  switch (event) {
  | Tick => handleTick(model)
  | _ => (model, [])
  };
};
