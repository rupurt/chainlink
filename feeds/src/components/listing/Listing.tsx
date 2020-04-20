import { DispatchBinding } from '@chainlink/ts-helpers'
import React, { useEffect } from 'react'
import { connect, MapStateToProps } from 'react-redux'
import { Row } from 'antd'
import GridItem from './GridItem'
import { AppState } from 'state'
import { listingOperations, listingSelectors } from '../../state/ducks/listing'

interface OwnProps {
  enableHealth: boolean
  compareOffchain: boolean
}

interface StateProps {
  loadingFeeds: boolean
  feedGroups: listingSelectors.ListingGroup[]
}

interface DispatchProps {
  fetchFeeds: DispatchBinding<typeof listingOperations.fetchFeeds>
  /* fetchHealthStatus: any */
}

interface Props extends OwnProps, StateProps, DispatchProps {}

export const Listing: React.FC<Props> = ({
  fetchFeeds,
  loadingFeeds,
  feedGroups,
  compareOffchain,
  enableHealth,
}) => {
  useEffect(() => {
    fetchFeeds()
  }, [fetchFeeds])
  /* useEffect(() => { */
  /*   if (enableHealth) { */
  /*     fetchHealthStatus(groups) */
  /*   } */
  /* }, [enableHealth, fetchHealthStatus, groups]) */

  let content
  if (loadingFeeds) {
    content = <>Loading Feeds...</>
  } else {
    content = (
      <div className="listing">
        {feedGroups.map(g => (
          <div className="listing-grid__group" key={g.name}>
            <h3 className="listing-grid__header">
              Decentralized Price Reference Data for {g.name} Pairs
            </h3>

            <Row gutter={18} className="listing-grid">
              {g.feeds.map(f => (
                <GridItem
                  key={f.name}
                  feed={f}
                  compareOffchain={compareOffchain}
                  enableHealth={enableHealth}
                />
              ))}
            </Row>
          </div>
        ))}
      </div>
    )
  }

  return content
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    loadingFeeds: state.listing.loadingFeeds,
    feedGroups: listingSelectors.feedGroups(state),
  }
}

const mapDispatchToProps = {
  fetchFeeds: listingOperations.fetchFeeds,
  /* fetchHealthStatus: listingOperations.fetchHealthStatus, */
}

export default connect(mapStateToProps, mapDispatchToProps)(Listing)
