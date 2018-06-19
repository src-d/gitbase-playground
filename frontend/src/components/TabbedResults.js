import React, { Component, Fragment } from 'react';
import { Row, Col, Alert, Tabs, Tab, Button, Glyphicon } from 'react-bootstrap';
import PropTypes from 'prop-types';
import ResultsTable from './ResultsTable';
import HistoryTable from './HistoryTable';
import Loader from './Loader';
import './TabbedResults.less';

class TabTitle extends Component {
  constructor(props) {
    super(props);

    this.state = {
      inEdit: false,
      title: props.title
    };
    this.ref = React.createRef();

    this.handleStartEdit = this.handleStartEdit.bind(this);
    this.handleEndEdit = this.handleEndEdit.bind(this);
  }

  handleStartEdit() {
    this.setState({ inEdit: true });
  }

  handleEndEdit() {
    this.setState({ inEdit: false });
  }

  render() {
    const { tabKey, active } = this.props;
    const { title, inEdit } = this.state;

    if (inEdit) {
      return (
        <div ref={this.ref}>
          <input
            type="text"
            className="tab-title"
            value={title}
            onChange={e => {
              this.setState({ title: e.target.value });
            }}
            onKeyPress={e => {
              if (e.key === 'Enter') {
                this.handleEndEdit();
              }
            }}
            onBlur={this.handleEndEdit}
          />
        </div>
      );
    }

    return (
      <div>
        <span className="tab-title">{title}</span>
        <Button
          className="btn-title"
          bsStyle={active ? 'gbpl-tertiary' : 'gbpl-primary-tint-2'}
          bsSize="xsmall"
          onClick={() => {
            this.handleStartEdit(tabKey);
          }}
        >
          <Glyphicon glyph="pencil" />
        </Button>
        <Button
          className="btn-title"
          bsStyle={active ? 'gbpl-tertiary' : 'gbpl-primary-tint-2'}
          bsSize="xsmall"
          onClick={() => {
            this.props.handleRemoveResult(tabKey);
          }}
        >
          <span aria-hidden="true">&times;</span>
        </Button>
      </div>
    );
  }
}

TabTitle.propTypes = {
  tabKey: PropTypes.any.isRequired,
  active: PropTypes.bool.isRequired,
  title: PropTypes.string.isRequired,
  handleRemoveResult: PropTypes.func.isRequired
};

class TabbedResults extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activeKey: 0,
      nTabs: 0
    };

    this.handleSelect = this.handleSelect.bind(this);
  }

  static getDerivedStateFromProps(nextProps, prevState) {
    const nextNTabs = nextProps.results.size;

    if (prevState.nTabs === nextNTabs) {
      return null;
    }

    // Make the last tab active when a new one is added,
    // or when the current active tab is deleted
    const newTab = prevState.nTabs < nextNTabs;
    const lostTab = !nextProps.results.has(prevState.activeKey);

    if (newTab || lostTab) {
      return {
        activeKey: Array.from(nextProps.results.keys())[nextNTabs - 1],
        nTabs: nextNTabs
      };
    }

    return { nTabs: nextNTabs };
  }

  handleSelect(activeKey) {
    this.setState({ activeKey });
    this.props.handleSetActiveResult(activeKey);
  }

  render() {
    const { showCode, showUAST, history } = this.props;

    return (
      <div className="results-padding full-height full-width">
        <Tabs
          id="tabbed-results"
          className="full-height"
          activeKey={this.state.activeKey}
          onSelect={this.handleSelect}
        >
          {Array.from(this.props.results.entries()).map(([key, query]) => {
            let content = '';
            if (key === this.state.activeKey) {
              if (query.loading) {
                content = (
                  <Row>
                    <Col className="text-center loader-col" xs={12}>
                      <Loader />
                    </Col>
                  </Row>
                );
              } else if (query.errorMsg) {
                content = (
                  <Row className="errors-row">
                    <Col xs={12}>
                      <Alert bsStyle="danger">{query.errorMsg}</Alert>
                    </Col>
                  </Row>
                );
              } else if (query.response) {
                content = (
                  <ResultsTable
                    response={query.response}
                    showCode={showCode}
                    showUAST={showUAST}
                  />
                );
              } else {
                content = (
                  <Fragment>
                    <Row>
                      <Col xs={12} className="text-center">
                        SUSPENDED TAB
                      </Col>
                    </Row>
                    <Row>
                      <Col xs={12} className="text-center">
                        <Button
                          className="reload"
                          bsStyle="gbpl-tertiary"
                          onClick={() => this.props.handleReload(key)}
                        >
                          Reload
                        </Button>
                      </Col>
                    </Row>
                  </Fragment>
                );
              }
            }

            return (
              <Tab
                key={key}
                eventKey={key}
                title={
                  <TabTitle
                    title={query.title || query.sql}
                    tabKey={key}
                    active={key === this.state.activeKey}
                    handleRemoveResult={this.props.handleRemoveResult}
                  />
                }
              >
                <Row className="query-row">
                  <Col xs={12}>
                    <div className="query-text">
                      <p>{query.sql}</p>
                    </div>
                    <Button
                      className="edit-query"
                      bsStyle="gbpl-tertiary-tint-2-link"
                      onClick={() => this.props.handleEditQuery(query.sql)}
                    >
                      EDIT
                    </Button>
                  </Col>
                </Row>
                {content}
              </Tab>
            );
          })}
          {history.length > 0 && (
            <Tab
              key="history"
              eventKey="history"
              tabClassName="history-tab-title"
              title={
                <div className="history-tab">
                  <span className="icon-bg">
                    <Glyphicon glyph="time" className="history-icon" />
                  </span>
                  <span className="tab-title">history</span>
                </div>
              }
            >
              <HistoryTable
                items={history}
                onOpenQuery={this.props.handleEditQuery}
                handleReset={this.props.handleResetHistory}
              />
            </Tab>
          )}
        </Tabs>
      </div>
    );
  }
}

TabbedResults.propTypes = {
  // results is a Map of objects, each object may contain:
  // sql: 'string'      Required
  // loading: true      Optional, tab will show a loading animation
  // errorMsg: 'string' Optional
  // response: object   Required if loading and errorMsg are not present
  results: PropTypes.instanceOf(Map).isRequired,
  history: HistoryTable.propTypes.items,
  handleRemoveResult: PropTypes.func.isRequired,
  handleEditQuery: PropTypes.func.isRequired,
  handleResetHistory: PropTypes.func.isRequired,
  handleSetActiveResult: PropTypes.func.isRequired,
  handleReload: PropTypes.func.isRequired,
  showCode: PropTypes.func.isRequired,
  showUAST: PropTypes.func.isRequired
};

export default TabbedResults;
