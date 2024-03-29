use std::fmt::Debug;
use std::cmp::PartialEq;
use std::result::Result;

#[derive(PartialEq, Debug, serde::Serialize, serde::Deserialize)]
pub enum Status {
  Paused,
  Running,
}


#[derive(PartialEq, Debug)]
pub enum Error {
  AlreadyExists,
  NotFound,
}

pub enum Command {
  CreateMonitoring(CreateMonitoring),
  PauseMonitoring(PauseMonitoring),
  ResumeMonitoring(ResumeMonitoring),
}

pub struct ResumeMonitoring {
  pub id: String,
}

pub struct CreateMonitoring {
  pub id: String,
  pub url: String,
}

pub struct PauseMonitoring {
  pub id: String,
}

#[derive(PartialEq, Debug, serde::Serialize)]
pub struct Monitoring {
  pub id: Option<String>,
  pub status: Status,
}


#[derive(PartialEq, Debug)]
pub enum Event {
  MonitoringStarted(MonitoringStarted),
  MonitoringPaused(MonitoringPaused),
  MonitoringResumed(MonitoringResumed),
}


#[derive(PartialEq, Debug, serde::Serialize, serde::Deserialize)]
pub struct MonitoringResumed {
  id: String,
}

#[derive(PartialEq, Debug, serde::Serialize, serde::Deserialize)]
pub struct MonitoringStarted {
  id: String,
  url: String,
}

#[derive(PartialEq, Debug, serde::Serialize, serde::Deserialize)]
pub struct MonitoringPaused {
  id: String,
}

pub fn initial_state() -> Monitoring {
  Monitoring { id: None, status: Status::Paused }
}

pub fn is_terminal(_state: &Monitoring) -> bool { false }

pub fn decide(state: &Monitoring, command: &Command) -> Result<Vec<Event>, Error> {
  match command {
    Command::CreateMonitoring(CreateMonitoring { id, url }) => {
      if state.id.is_some() {
        return Err(Error::AlreadyExists);
      }

      Ok(vec![
        Event::MonitoringStarted(MonitoringStarted { id: id.to_string(), url: url.to_string() })
      ])
    }
    Command::PauseMonitoring(PauseMonitoring { id }) => {
      if state.id.is_none() {
        return Err(Error::NotFound);
      }

      Ok(vec![Event::MonitoringPaused(MonitoringPaused { id: id.to_string() })])
    }
    Command::ResumeMonitoring(ResumeMonitoring { id }) => {
      if state.id.is_none() {
        return Err(Error::NotFound);
      }

      Ok(vec![Event::MonitoringResumed(MonitoringResumed { id: id.to_string() })])
    }
  }
}

pub fn evolve(state: &Monitoring, event: &Event) -> Monitoring {
  match event {
    Event::MonitoringStarted(MonitoringStarted{ id, .. }) => {
      Monitoring { id: Some(id.to_string()), status: Status::Running }
    }
    Event::MonitoringPaused(MonitoringPaused{ .. }) => {
      Monitoring { status: Status::Paused, id: state.id.clone() }
    }
    Event::MonitoringResumed(MonitoringResumed{ .. }) => {
      Monitoring { status: Status::Running, id: state.id.clone() }
    }
  }
}

#[cfg(test)]
mod tests {
  use super::*;
  use infra::decider;
  use infra::testing;

  #[test]
  fn it_works() {
    let monitoring = decider::Decider::new(
      decide,
      evolve,
      initial_state,
      Some(is_terminal),
    );

    testing::Spec::new(monitoring)
      .given(vec![])
      .when(&Command::CreateMonitoring(CreateMonitoring {
        id: String::from("1"),
        url: String::from("https://example.com"),
      })
      )
      .then(testing::SpecResult::Event {
        events: vec![
          Event::MonitoringStarted {
            id: String::from("1"),
            url: String::from("https://example.com"),
          }
        ]
      });
  }
}
