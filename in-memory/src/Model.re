type configuration = {cozyURL: string};

type model = {
  cozyURL: string,
  ticked: int,
};

let init = (config: configuration): model => {
  Js.log(("init", config));
  {cozyURL: config.cozyURL, ticked: 0};
};

let update = (current: model, event: Event.event): (model, unit) => {
  Js.log(("update", current, event));
  let next =
    switch (event) {
    | Tick => {...current, ticked: current.ticked + 1}
    | _ => current
    };
  (next, ());
};
