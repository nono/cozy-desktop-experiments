type configuration = {cozyURL: string};

type ticks = int;

type model = {
  cozyURL: string,
  ticked: ticks,
};

let init = (config: configuration): model => {
  {cozyURL: config.cozyURL, ticked: 0};
};

let update = (current: model, event: Event.event): (model, unit) => {
  let next =
    switch (event) {
    | Tick => {...current, ticked: current.ticked + 1}
    | _ => current
    };
  (next, ());
};
