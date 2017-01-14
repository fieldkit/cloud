
import React, {PropTypes} from 'react'

const AboutPage = () => {
  return (
    <div className='page'  id="aboutPage">

      <div className="pageWrapper">
        <iframe className="vimeo" src="https://player.vimeo.com/video/124421450?autoplay=0&api=1" width={window.innerWidth * 0.9} height={window.innerWidth * 0.9 * 0.525} frameBorder="0"allowFullScreen></iframe>
        <div className="columnWrapper">
          <div className="column headline">
            <p>
              18 days, 345 kilometers,<br/>
              1 river, 31 adventurers, 100% open data.<br/>
              Join us in real-time as we explore
            </p>
            <h1>
              THE BEATING HEART OF OUR PLANET
            </h1>
          </div>
          <div className="column">
            <p>
              The Okavango Delta is one of the world’s last great wetland wildernesses. Although the Delta has been awarded UNESCO WHS Status its catchments in the highlands of Angola are still unprotected and largely unstudied. A team of Ba’Yei, scientists, engineers and adventurers will journey a 345 kilometers crossing the river, finding new species, exploring new ground, and taking the pulse of this mighty river that brings life-giving water to the Jewel of the Kalahari.
            </p>
            <p>
              This site displays data which is uploaded daily, via satellite, by the expedition team. Data is also available through a public API, allowing anyone to remix, analyze or visualize the collected information.
            </p>
          </div>
        </div>
        <div className="columnWrapper">
          <div className="column">
            <div class="goalIcon"><img src="/static/img/iconIntroUnderstand.png"/></div>
            <h2>UNDERSTAND<br/>THE WILDERNESS</h2>
            <p>
              To effectively protect the Okavango and its catchments it is essential to gain knowledge and insight into the functioning of the system as a whole. Starting in 2011 the Okavango Wilderness Project has conducted yearly transects of the Delta, gathering unique data and immersing the expedition members in the ebb and flow of this pristine wilderness. Traveling down the river, the team will collect data on insects, fish, birds, reptiles and mammals, as well as conduct water quality assessments and landscape surveys.
            </p>
          </div>
          <div className="column">
            <div class="goalIcon"><img src="/static/img/iconIntroPreserve.png"/></div>
            <h2>RAISE AWARENESS<br/>AND PRESERVE</h2>
            <p>
              Although the Okavango itself is protected as a UNESCO World Heritage Site, its catchment and water supply in Angola and Namibia remain vulnerable to human interference. By gathering and freely disseminating information about the functioning and health of the entire system the 2016 expedition aims to raise the levels of both understanding and awareness of this unique and fragile system.
            </p>
            <p>
              Once base-line data on the system becomes freely available effective measures can then be implemented to insure the continued health and survival of this great African wilderness.
            </p>
          </div>
        </div>
        <div className="columnWrapper credits">
          <div className="column">
            <h2>
              EXPEDITION TEAM<br/><span class="job"><span class="explorerBox legend"></span> National Geographic Emerging Explorers</span>
            </h2>
            <p>
              Adjany Costa <span class="job">Assistant Director & 2nd Fish</span><br/>
              Chris Boyes <span class="job">Expeditions Leader</span><br/>
              Gobonamang "GB" Kgetho <span class="job">Poler</span><br/>
              Gotz Neef <span class="job">Scienetific Collections & Leader Invertebrates</span><br/>
              Jer Thorp <span class="job">Data</span><br/>
              John Hilton <span class="job">Project Director</span><br/>
              Kerllen Costa <span class="job">Plants & Mammals</span><br/>
              Kyle Gordon <span class="job">Expedition Logistics</span><br/>
              Leilamang "Schnapps" Kgetho <span class="job">Poler</span><br/>
              Luke Manson <span class="job">Expedition Logistics</span><br/>
              Mia Maestro <span class="job">Photographer</span><br/>
              Motiemang “Judge” Xhikabora <span class="job">Poler</span><br/>
              Neil Gelinas <span class="job">Filmmaker</span><br/>
              Ninda Baptista <span class="job">Reptiles</span><br/>
              Nkeletsang “Ralph” Moshupa <span class="job">Poler</span><br/>
              Rachel Sussman <span class="job">Photographer</span><br/>
              Shah Selbe <span class="job">Tech</span><br/>
              Steve Boyes <span class="job">Project Leader & Birds</span><br/>
              Topho "Tom" Retiyo <span class="job">Poler</span><br/>
              Tumeleto "Water" Setlabosha<span class="job">Poler</span><br/>
            </p>
          </div>
          <div className="column">
            <div class="logos">
              <a href="http://www.nationalgeographic.com/">
                <img src="/static/img/natgeoLogo.svg" alt="National Geographic Logo" height="35px"/>
              </a>
              <a href="http://conservify.org/">
                <img src="/static/img/conservify.png" alt="Conservify Logo" height="35px"/>
              </a>
              <a href="http://www.o-c-r.org/">
                <img src="/static/img/ocrLogo.svg" alt="The Office for Creative Research Logo" height="35px"/>
              </a>
              <a href="http://www.wildbirdtrust.com/">
                <img src="/static/img/wbtLogo.png" alt="Wild Bird Trust Logo" height="35px"/>
              </a>
            </div>
          </div>
        </div>
      </div>

    </div>
  )
}

AboutPage.propTypes = {
  // active: PropTypes.bool.isRequired
}

export default AboutPage
