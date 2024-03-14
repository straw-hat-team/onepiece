use crate::decider;
use std::fmt::Debug;
use std::iter::Map;

pub type StepDecideReducer<T, State, Command, Event, Error> = fn(
    state: &State,
    command: &Command,
    current_value: T,
    changes: &StepStates<State>,
) -> Result<Vec<Event>, Error>;

pub struct StepStates<State>(Map<String, State>);

impl <State> StepStates<State> {
    pub fn insert(&mut self, key: String, value: State) {
        self.0.insert(key, value);
    }
}

struct StepDecide<T, State, Command, Event, Error> {
    current_value: T,
    decide: StepDecideReducer<T, State, Command, Event, Error>,
}

enum Step<T, State, Command, Event, Error> {
    DecideReducer(StepDecide<T, State, Command, Event, Error>),
    Decide(decider::Decide<State, Command, Event, Error>),
}

struct MultiStep<T, State, Command, Event, Error> {
    step_name: String,
    decide: Step<T, State, Command, Event, Error>,
}

pub struct Multi<T, State, Command, Event, Error> {
    state: State,
    evolve: decider::Evolve<State, Event>,
    steps: Vec<MultiStep<T, State, Command, Event, Error>>,
}

impl<T, State, Command, Event, Error> Multi<T, State, Command, Event, Error>
where
    Event: PartialEq + Debug,
    Error: PartialEq + Debug,
    State: PartialEq + Debug,
{
    pub fn new(evolve: decider::Evolve<State, Event>, state: State) -> Self {
        Multi {
            evolve,
            state,
            ..Default::default()
        }
    }

    pub fn execute(
        &mut self,
        step_name: String,
        decide: Step::Decide<T, State, Command, Event, Error>,
    ) -> &mut Self {
        self.steps.push(MultiStep { step_name, decide });
        self
    }

    pub fn reduce<T>(
        &mut self,
        step_name: String,
        lists: Vec<T>,
        decide: StepDecideReducer<T, State, Command, Event, Error>,
    ) -> &mut Self {
        for value in lists {
            self.steps.push(MultiStep {
                step_name: step_name.clone(),
                decide: Step::DecideReducer(StepDecide {
                    current_value: value,
                    decide,
                }),
            });
        }
        self
    }

    pub fn run(&self, command: &Command) -> Result<Vec<Event>, Error> {
        let mut acc = vec![];
        let mut current_state = self.state.clone();
        let mut step_states: StepStates<State> = StepStates(());

        for step in self.steps {
            let events = match step.decide {
                Step::Decide(decider) => decider(&self.state, command)?,
                Step::DecideReducer(decider) => decider.decide(
                    current_state,
                    command,
                    decider.current_value,
                    &step_states,
                )?,
            };

            current_state = events.iter().fold(current_state, self.evolve);
            step_states.insert(step.step_name, current_state.clone());
            acc.extend(events);
        }

        Ok(acc)
    }
}
