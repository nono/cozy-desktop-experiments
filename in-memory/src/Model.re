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

type model = {
  config: configuration,
  ticked: ticks,
  changesState: changesFeedState,
};

let init = (config: configuration): model => {
  {
    config: config,
    ticked: 0,
    changesState: ChangesFeedNeverFetched,
  };
};
