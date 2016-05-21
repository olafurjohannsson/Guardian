using System;
using System.Collections.Generic;
using System.Collections.Concurrent;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Threading;

namespace Guardian
{
    public class Producer<T> : IProducer<T>
    {
        private BlockingCollection<T> collection;
        private TimeSpan timeout;
        private CancellationToken cancellationToken;

        public Producer(BlockingCollection<T> collection)
        {
            this.collection = collection;
            this.timeout = TimeSpan.FromSeconds(2);
            this.cancellationToken = new CancellationToken();
            this.cancellationToken.Register(() =>
            {
                // .... ?
            });
        }

        public bool Produce(T item)
        {
            if (this.collection == null)
            {
                throw new NullReferenceException(string.Format("Collection is null in {0}", this.GetType().Name));
            }

            try
            {
                while (true)
                {

                    // Well we don't want to add if someone cancelled now do we.
                    if (this.cancellationToken != null && this.cancellationToken.IsCancellationRequested)
                    {
                        return false;
                    }

                    return this.collection.TryAdd(item, (int)this.timeout.TotalMilliseconds);
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.ToString());
            }

            return false;
        }
    }
}
