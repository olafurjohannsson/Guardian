using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Guardian
{

    public class GuardianEntity
    {
        private DateTime now;
        private HashSet<GuardianRequest> requests;

        public GuardianEntity()
        {
            this.now = DateTime.Now;
            this.requests = new HashSet<GuardianRequest>();
        }

        public IEnumerable<GuardianRequestGroupEntity> GetRequestsToday()
        {
            // Start by filtering requests today
            var requestsByDay = from request in this.requests
                                where request.TimeStamp.DayOfYear == this.now.DayOfYear
                                select request;

            // group hosts together and return group
            return from request in requestsByDay
                   group request by request.Host into g
                   select new GuardianRequestGroupEntity
                   {
                       Host = g.Key,
                       Started = g.Min(x => x.TimeStamp),
                       Ended = g.Max(x => x.TimeStamp),
                       Count = g.Count()
                   };
        }

        public int Count()
        {
            return this.requests.Count;
        }

        public bool Add(GuardianRequest request)
        {
            return this.requests.Add(request);
        }
    }
}
