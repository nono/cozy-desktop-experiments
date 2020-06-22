type effect =
  | FetchChangesFeed;

type effects = list(effect);

type event =
  | Tick
  | ReceiveChangesFeed(Remote.changes);

type configuration = {cozyURL: string};

type ticks = int;

type changesFeedStatus =
  | ChangesFeedNeverFetched
  | ChangesFeedCurrentlyFetching
  | ChangesFeedLastFetchedAt(ticks);

type model = {
  cozyURL: string,
  ticked: ticks,
  changesStatus: changesFeedStatus,
};

let init = (config: configuration): model => {
  {
    cozyURL: config.cozyURL,
    ticked: 0,
    changesStatus: ChangesFeedNeverFetched,
  };
};

type handler = model => (model, effects);

let incrementTicked: handler =
  current => {
    ({...current, ticked: current.ticked + 1}, []);
  };

let checkFetchChanges: handler =
  current => {
    switch (current.changesStatus) {
    | ChangesFeedNeverFetched => (
        {...current, changesStatus: ChangesFeedCurrentlyFetching},
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
