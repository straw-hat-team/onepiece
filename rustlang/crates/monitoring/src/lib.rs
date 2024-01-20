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


#[derive(PartialEq, Debug, serde::Serialize, serde::Deserialize)]
pub enum Event {
  MonitoringStarted {
    id: String,
    url: String,
  },
  MonitoringPaused {
    id: String,
  },
  MonitoringResumed {
    id: String,
  },
}

pub fn initial_state() -> Monitoring {
  Monitoring { id: None, status: Status::Paused }
}

pub fn is_terminal(_aggregate: &Monitoring) -> bool { false }

pub fn decide(aggregate: &Monitoring, command: &Command) -> Result<Vec<Event>, Error> {
  match command {
    Command::CreateMonitoring(CreateMonitoring { id, url }) => {
      if aggregate.id.is_some() {
        return Err(Error::AlreadyExists);
      }

      let event = Event::MonitoringStarted { id: id.to_string(), url: url.to_string() };
      Ok(vec![event])
    }
    Command::PauseMonitoring(PauseMonitoring { id }) => {
      if aggregate.id.is_none() {
        return Err(Error::NotFound);
      }

      let event = Event::MonitoringPaused { id: id.to_string() };
      Ok(vec![event])
    }
    Command::ResumeMonitoring(ResumeMonitoring { id }) => {
      if aggregate.id.is_none() {
        return Err(Error::NotFound);
      }

      let event = Event::MonitoringResumed { id: id.to_string() };
      Ok(vec![event])
    }
  }
}

pub fn evolve(aggregate: &Monitoring, event: &Event) -> Monitoring {
  match event {
    Event::MonitoringStarted { id, .. } => {
      Monitoring { id: Some(id.to_string()), status: Status::Running }
    }
    Event::MonitoringPaused { .. } => {
      Monitoring { status: Status::Paused, id: aggregate.id.clone() }
    }
    Event::MonitoringResumed { .. } => {
      Monitoring { status: Status::Running, id: aggregate.id.clone() }
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
