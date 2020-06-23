type effect =
  | FetchChangesFeed;

type effects = list(effect);

type event =
  | Tick
  | ReceiveChangesFeed(Remote.changes);

type config = {cozyURL: string};

type ticks = int;

type changesFeedState =
  | ChangesFeedNeverFetched
  | ChangesFeedCurrentlyFetching
  | ChangesFeedLastFetchedAt(ticks);

type states = {changes: changesFeedState};

type t = {
  config,
  ticked: ticks,
  states,
  remote: Remote.tree,
};

let init = (config: config): t => {
  {
    config,
    ticked: 0,
    states: {
      changes: ChangesFeedNeverFetched,
    },
    remote: Remote.emptyTree,
  };
};
