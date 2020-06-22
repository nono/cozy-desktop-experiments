type effect =
  | FetchChangesFeed;

type effects = list(effect);

type event =
  | Tick
  | ReceiveChangesFeed(Remote.changes);

type configuration = {cozyURL: string};

type ticks = int;

type changesFeedState =
  | ChangesFeedNeverFetched
  | ChangesFeedCurrentlyFetching
  | ChangesFeedLastFetchedAt(ticks);

type states = {changes: changesFeedState};

type t = {
  config: configuration,
  ticked: ticks,
  states,
  remote: Remote.tree,
};

let init = (config: configuration): t => {
  {
    config,
    ticked: 0,
    states: {
      changes: ChangesFeedNeverFetched,
    },
    remote: Remote.emptyTree,
  };
};
